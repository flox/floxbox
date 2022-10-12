{ mkShell, qemu, qemu-utils, ... }:

mkShell {
  buildInputs = [ qemu qemu-utils ];
}
