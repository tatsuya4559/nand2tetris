open Core
open Jackc

let die msg =
    Printf.eprintf "%s\n" msg; exit 1

let out_filename arg =
    let (stem, _) = Filename.split_extension arg in
    stem ^ ".vm"

let compile_file filename =
    In_channel.with_file filename ~f:(fun ic ->
        let src = In_channel.input_all ic in
        match Parser.parse Parser.parse_class src with
        | None -> Printf.eprintf "error!"; exit 1
        | Some (parsed, _) ->
            Out_channel.write_all (out_filename filename) ~data:(Ast.show_class_dec parsed)
    )

let compile_dir dirname =
    Sys_unix.ls_dir dirname
    |> List.filter ~f:(fun f -> Filename.check_suffix f ".jack")
    |> List.iter ~f:compile_file

let compile arg =
    if Poly.equal (Sys_unix.is_file arg) `Yes then
        compile_file arg
    else if Poly.equal (Sys_unix.is_directory arg) `Yes then
        compile_dir arg
    else
        die @@ (Printf.sprintf "%s is not filename or dirname" arg)

let command =
    Command.basic
    ~summary:"Jack compiler"
    Command.Param.(
        map
        (anon ("filename or dirname" %: string))
        ~f:(fun file_or_dir () -> compile file_or_dir))

let () =
    Command_unix.run command
