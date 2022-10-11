## Currently only supported running on Linux host machines with nix installed. Currently only supports `ubuntu-focal` vm. More VM's to come!
## Bootstrapping sequence

1. `nix profile install --refresh github:flox/flox-qemu`
2. run `flox-qemu init-image` (wait for install to complete, if installer crashes, re-run `flox-qemu init-image`)
3. Wait until the installation is complete, hit `<enter>` and login with user: flox-qemu password: flox-qemu, then type `poweroff` and `<enter>` to complete the install process.
4. run `flox-qemu images-list` and copy the BASE image you want to snapshot.
5. run `flox-qemu snapshot-image --base-image-name=2022-07-27T14-00-09-605c6a8e-6ce4-489f-9ed0-598e76bd0d31-flox-qemu-BASE-ubuntu.img.qcow2 --snapshot-name=mysnapshot` (replace base-image-name with your actual image name)
6. run `flox-qemu run-image --image-name=2022-07-27T14-05-58-flox-qemu-SNAPSHOT-ubuntu.img.qcow2-mysnapshot` (replace `--image-name` with the name of the snapshot you want to run).
7. Use `Ctr+C` to stop the image from running.


Now you can ssh into this machine at ssh flox-qemu@localhost with password `flox-qemu`
