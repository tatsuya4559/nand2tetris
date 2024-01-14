open Core
open Ast

(*** Parsec ***)
type 'a t = char list -> ('a * char list) option

let parse p s = p @@ String.to_list s

let get_char = function
    | [] -> None
    | c :: cs -> Some (c, cs)

let map f p =
    fun cs -> match p cs with
    | None -> None
    | Some (c, cs') -> Some (f c, cs')
let ( <$> ) = map

let pure v = fun cs -> Some (v, cs)

let apply fp xp =
    fun cs -> match fp cs with
    | None -> None
    | Some (f, cs') -> (map f xp) cs'

let ( <*> ) = apply
let ( <* ) xp yp = (fun x _ -> x) <$> xp <*> yp
let ( *> ) xp yp = (fun _ y -> y) <$> xp <*> yp

let product xp yp = (fun x y -> (x, y)) <$> xp <*> yp
let ( <.> ) xa ya = product xa ya

let empty = fun _ -> None

let either xp yp =
    fun cs -> match xp cs with
    | None -> yp cs
    | Some _ as result -> result

let ( <|> ) = either

let bind xp fp =
    fun cs -> match xp cs with
    | None -> None
    | Some (x, cs') -> fp x cs'

let ( >>= ) = bind
let ( let* ) xp fp = bind xp fp
let return = pure

let satisfy predicate =
    let* x = get_char in
    if predicate x then
        pure x
    else
        empty

let match_char c = satisfy (Char.equal c)

let match_space =
    let p = function
        | ' ' | '\t' | '\n' | '\r' -> true
        | _ -> false
    in
    satisfy p

let rec many p cs =
    (some p <|> pure []) cs
and some p cs =
    (List.cons <$> p <*> many p) cs

let optional default p = p <|> pure default

(*** Matcher ***)
let match_digit =
    let p = function
        | '0' .. '9' -> true
        | _ -> false
    in
    satisfy p

let match_letter =
    let p = function
        | 'a' .. 'z' | 'A' .. 'Z' | '_' -> true
        | _ -> false
    in
    satisfy p

let match_alnum = match_letter <|> match_digit

let match_ident =
    let match_ident_list = List.cons <$> match_letter <*> many match_alnum in
    String.of_char_list <$> match_ident_list

let match_int =
    int_of_string <$> (String.of_char_list <$> some match_digit)

let match_non_quote =
    satisfy (fun c -> Char.(c <> '"'))

let match_string =
    let match_string_list =  match_char '"' *> many match_non_quote <* match_char '"' in
    String.of_char_list <$> match_string_list

let match_symbol s =
    let rec f = function
        | [] -> pure []
        | c :: cs ->
            let* _ = match_char c in
            let* _ = f cs in
            pure (c :: cs)
    in
    String.to_list s |> f |> map String.of_char_list

let match_true =
    let* _ = match_symbol "true" in
    pure true

let match_false =
    let* _ = match_symbol "false" in
    pure false

let match_bool = match_true <|> match_false


(*** Getter ***)
let get_token p = many match_space *> p <* many match_space

let get_ident = get_token match_ident
let get_int = get_token match_int
let get_string = get_token match_string
let get_symbol s = get_token @@ match_symbol s
let get_bool = get_token match_bool

let separated ch p = List.cons <$> p <*> (many (get_symbol ch *> p))

(*** Parser ***)
let just sym value = (fun _ -> value) <$> get_symbol sym

(* ident *)
let parse_ident = (fun x -> Ident x) <$> get_ident

(* constant *)
let parse_int = (fun x -> Int x) <$> get_int
let parse_string = (fun x -> String x) <$> get_string

(* keyword *)
let parse_bool = (fun x -> Bool x) <$> get_bool
let parse_null = just "null" Null
let parse_this = just "this" This
let parse_keyword = parse_bool <|> parse_null <|> parse_this

(* unary op *)
let parse_neg = just "-" Neg
let parse_not = just "~" Not
let parse_unary_op = parse_neg <|> parse_not

(* binary op *)
let parse_plus = just "+" Plus
let parse_minus = just "-" Minus
let parse_mul = just "*" Mul
let parse_div = just "/" Div
let parse_and = just "&" And
let parse_or = just "|" Or
let parse_lt = just "<" Lt
let parse_gt = just ">" Gt
let parse_eq = just "=" Eq
let parse_binary_op =
    parse_plus <|> parse_minus <|> parse_mul <|> parse_div <|>
    parse_and <|> parse_or <|>
    parse_lt <|> parse_gt <|> parse_eq

(* expr *)
let rec parse_term cs =
    (parse_int <|> parse_string <|> parse_keyword <|>
    parse_indexing <|> parse_call <|> parse_ident <|>
    parse_prefix <|> parse_factor) cs

and parse_expr cs = (
    let* left = parse_term in
    let* result = many (parse_binary_op <.> parse_expr) in
    match result with
    | [] -> return left
    | (op, right) :: tl ->
        let init = Infix (left, op, right) in
        return @@ List.fold ~init tl ~f:(fun acc (o, r) -> Infix (acc, o, r))
    ) cs

and parse_prefix cs =
    ((fun (op, e) -> Prefix (op, e)) <$> (parse_unary_op <.> parse_expr)) cs

and parse_factor cs =
    (get_symbol "(" *> parse_expr <* get_symbol ")") cs

and parse_indexing cs = (
    let* e = parse_ident in
    let* idx = get_symbol "[" *> parse_expr <* get_symbol "]" in
    return @@ Indexing (e, idx)
    ) cs

and parse_call cs = (
    let parse_args =
        get_symbol "(" *>
        optional [] (separated "," parse_expr)
        <* get_symbol ")"
    in
    let parse_dot_notation =
        let* left = parse_ident in
        let* right = get_symbol "." *> parse_ident in
        return @@ Dot_notation (left, right)
    in
    let* fn = parse_dot_notation <|> parse_ident in
    let* args = parse_args in
    return @@ Call (fn, args)
    ) cs

(* statement *)
let parse_do =
    let* call = get_symbol "do" *> parse_call <* get_symbol ";" in
    return @@ Do call

let parse_return =
    let* v = get_symbol "return" *> optional None ((fun x -> Some x) <$> parse_expr) <* get_symbol ";" in
    return @@ Return v

let parse_let =
    let* name = get_symbol "let" *> parse_term in
    let* v = get_symbol "=" *> parse_expr <* get_symbol ";" in
    return @@ Let (name, v)

let rec parse_statement cs =
    (parse_do <|> parse_return <|> parse_let <|> parse_if <|> parse_while) cs

and parse_while cs = (
    let* cond = get_symbol "while" *> get_symbol "(" *> parse_expr <* get_symbol ")" in
    let* body = get_symbol "{" *> many parse_statement <* get_symbol "}" in
    return @@ While { cond; body }
    ) cs

and parse_if cs = (
    let* cond = get_symbol "if" *> get_symbol "(" *> parse_expr <* get_symbol ")" in
    let* then_clause = get_symbol "{" *> many parse_statement <* get_symbol "}" in
    let* else_clause = optional [] (get_symbol "else" *> get_symbol "{" *> many parse_statement <* get_symbol "}") in
    return @@ If { cond; then_clause; else_clause }
    ) cs

(* declaration *)

let parse_field = just "field" Field
let parse_static = just "static" Static
let parse_storage  = parse_field <|> parse_static

let parse_int_type = just "int" Int_type
let parse_bool_type = just "boolean" Bool_type
let parse_char_type = just "char" Char_type
let parse_class_name = get_ident >>= (fun name -> return @@ Class_name name)
let parse_type = parse_int_type <|> parse_bool_type <|> parse_char_type <|> parse_class_name

let parse_void = just "void" Void
let parse_return_type = parse_void <|> parse_type


let parse_local_var =
    let* typ = get_symbol "var" *> parse_type in
    let* names = separated "," get_ident <* get_symbol ";" in
    return { typ; names }

let parse_class_var =
    let* storage = parse_storage in
    let* typ = parse_type in
    let* names = separated "," get_ident <* get_symbol ";" in
    return { storage; typ; names }

let parse_constructor = just "constructor" Constructor
let parse_function = just "function" Function
let parse_method = just "method" Method
let parse_subroutine_kind = parse_constructor <|> parse_function <|> parse_method

let parse_subroutine_param =
    let* typ = parse_type in
    let* name = get_ident in
    return { typ; name }

let parse_subroutine_body =
    let* vars = many parse_local_var in
    let* stmts = many parse_statement in
    return { vars; stmts }

let parse_subroutine =
    let* kind = parse_subroutine_kind in
    let* return_type = parse_return_type in
    let* name = get_ident in
    let* params = get_symbol "(" *> optional [] (separated "," parse_subroutine_param) <* get_symbol ")" in
    let* body = get_symbol "{" *> parse_subroutine_body <* get_symbol "}" in
    return { kind; return_type; name; params; body }

let parse_class =
    let* name = get_symbol "class" *> get_ident <* get_symbol "{" in
    let* vars = many parse_class_var in
    let* subroutines = many parse_subroutine <* get_symbol "}" in
    return {name; vars; subroutines }
