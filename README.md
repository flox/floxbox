# floxbox

`floxbox` is a box for your flox, that'll knock off your socks. You put flox into the box, and you've got a floxbox. 


# Deprecated

## Currently only supported running on Linux host machines with nix installed. Currently only supports `ubuntu-focal` vm. More VM's to come!
## Bootstrapping sequence

1. `nix profile install --refresh github:flox/floxbox`
2. run `floxbox init-image` (wait for install to complete, if installer crashes, re-run `floxbox init-image`)
3. Wait until the installation is complete, hit `<enter>` and login with user: floxbox password: floxbox, then type `poweroff` and `<enter>` to complete the install process.
4. run `floxbox images-list` and copy the BASE image you want to snapshot.
5. run `floxbox snapshot-image --base-image-name=2022-07-27T14-00-09-605c6a8e-6ce4-489f-9ed0-598e76bd0d31-floxbox-BASE-ubuntu.img.qcow2 --snapshot-name=mysnapshot` (replace base-image-name with your actual image name)
6. run `floxbox run-image --image-name=2022-07-27T14-05-58-floxbox-SNAPSHOT-ubuntu.img.qcow2-mysnapshot` (replace `--image-name` with the name of the snapshot you want to run).
7. Use `Ctr+C` to stop the image from running.


Now you can ssh into this machine at ssh floxbox@localhost with password `floxbox`
