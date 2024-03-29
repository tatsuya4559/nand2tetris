type unary_op =
    | Neg
    | Not
    [@@deriving show]

type binary_op =
    | Plus | Minus | Mul | Div
    | And | Or
    | Lt | Gt | Eq
    [@@deriving show]

type expr =
    | Ident of string
    | Int of int
    | String of string
    | Bool of bool
    | Null
    | This
    | Prefix of (unary_op * expr)
    | Infix of (expr * binary_op * expr)
    | Indexing of (expr * expr)
    | Dot_notation of (expr * expr)
    | Call of (expr * expr list)
    [@@deriving show]

type statement =
    | Let of (expr * expr)
    | If of { cond: expr; then_clause: statement list; else_clause: statement list }
    | While of { cond: expr; body: statement list }
    | Do of expr
    | Return of expr option
    [@@deriving show]

type storage_class =
    | Field
    | Static
    | Var
    [@@deriving show]

type jack_type =
    | Int_type
    | Bool_type
    | Char_type
    | Class_name of string
    | Void
    [@@deriving show]

type var_dec = {
    storage: storage_class;
    typ: jack_type;
    name: string;
}
[@@deriving show]

type subroutine_kind =
    | Constructor
    | Function
    | Method
    [@@deriving show]

type subroutine_param = {
    typ: jack_type;
    name: string
}
[@@deriving show]

type subroutine_body = {
    vars: var_dec list;
    stmts: statement list;
}
[@@deriving show]

type subroutine_dec = {
    kind: subroutine_kind;
    return_type: jack_type;
    name: string;
    params: subroutine_param list;
    body: subroutine_body;
}
[@@deriving show]

type class_dec = {
    name: string;
    vars: var_dec list;
    subroutines: subroutine_dec list;
}
[@@deriving show]
