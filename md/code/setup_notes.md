## Setup notes and memory dump

### `dwm` and `st`

- make changes to `config.def.h`, **not** `config.h`
- customisation options go before `exec dwm` in your `~/.xinitrc`
- terminal.sexy + st = nice st, but their default export uses _static_ unsigned
  ints instead of just ints for foreground and background colour. Change to
unsigned int.

### X + nvidia drivers

- legacy Nvidia driver I'm using completely fucks with X, just random freezing
  at boot. Seems to be resolved by creating a file called
`/etc/modprobe.d/modprobe.conf` containing

	options nvidia NVreg_Mobile=1

### Lynx

- the confirmation button in the options menu is **at the bottom**
- you \[a\]dd and \[v\]iew bookmarks
- you \[d\]ownload pages
