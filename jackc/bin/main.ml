open Core
open Jackc

let compile filename =
    let ic = In_channel.create filename in
    let src = In_channel.input_all ic in
    match Parser.parse Parser.parse_class src with
    | None -> Printf.eprintf "error!"; exit 1
    | Some (parsed, _) -> Ast.show_class_dec parsed |> print_endline
;;

let command =
    Command.basic
    ~summary:"Jack compiler"
    Command.Param.(
        map
        (anon ("filename" %: string))
        ~f:(fun filename () -> compile filename))

let () =
    Command_unix.run command
