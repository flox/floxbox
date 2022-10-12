{ stdenv, qemu ... }:

stdenv.mkDerivation {
  buildInputs = [ qemu qemu-utils ];
}
