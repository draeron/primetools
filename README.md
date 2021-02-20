# primetools

A small tools for synchronising meta data between
[ITunes](https://www.apple.com/ca/itunes/) and [Denon's Prime
Engine](https://www.denondj.com/engine-prime).

This tool came with my necessity to import particular meta data from ITunes to
PRIME Engine, mainly time added / modified. Also due since I lost my Engine
PRIME database, I needed a way to restore crates/playlists saved on a exported
disk back into my desktop database, a functionality which should be there but
doesn't exists.

## Disclaimer

_This goes without saying that you should use this at your own risk. This tool
could easily wipeout any of the database it connects to. Make some backup prior
to using it._

For that reason, I don't plan on providing compiled binary.

## Compiling

You will require both the [golang toolchain](https://golang.org/dl/) and a CGO
compatible C compiler such as [GCC](https://jmeubank.github.io/tdm-gcc/) because
of the dependency on sqlite.

The code should compile on non-windows platforms but it's been only developped
on that OS.

```bash
git clone https://github.com/draeron/primetools.git
cd primetools
go build
```

Et voilà! You should have a single binary called `primetools.exe`.

## Feature Matrix

### Sources

| Source    | Files | ITunes | PRIME | Rekordbox |
| --------- | ----- | ------ | ----- | --------- |
| Rating    | [x]   | [x]    | [x]   | [x]       |
| Playlists |       | [x]    | [x]   | [x]       |
| Crates    |       |        | [x]   | [x]       |
| Time      | [x]   | [x]    | [x]   | [x]       |

### Targets

| Target          | Files | ITunes\* | PRIME | Traktor |
| --------------- | ----- | -------- | ----- | ------- |
| Add Files       |       | [x]      | [ ]   |         |
| Fix Renames     |       | [x]      | [x]   |         |
| Fix Duplicate   |       | [ ]      | [ ]   |         |
| Sync Rating     | [x]   | [x]      | [x]   | [x]     |
| Sync Time       | [x]   | [x]      | [x]   |         |
| Dump Crates     |       |          | [x]   |         |
| Dump Playlist   |       | [x]      | [x]   |         |
| Import Crates   | [ ]   |          | [x]   |         |
| Import Playlist |       | [ ]      | [x]   |         |

_Legend_

    [x]   Implemented
    [ ]   Need Implementation
    empty Not Applicable


#### _MacOS_

Writing to Itunes is done implemented through windows COM interface. On MacOS
this require going through sytstem javascript scripting which I haven't
explored. Since I use only Windows I haven't implemented a writer for MacOS. You
can still use the tools to import mort of the meta data through the .plist file.

#### _Traktor_

Traktor is only supported because I can read the proper POPM id3 frame (ie:
rating) that is used by Traktor. Meta data from NML is not implemented.

#### What about _Serato_ ?

Since I don't really use Rekorbox and Serato but the code is modular enough that
they could be eventually implemented. If someone submit a pull request for it I
would gladly review and merge it.

#### Known Issues

1. My code doesn't likes slash character in playlist / crate names since that's
   the character used for expressing folder structure.

2. There might be collisions/issues with unicode characters, the reason being
   Itunes and PRIME doesn't store strings in the same unicode format so I had to
   do some approximation based on nearest non-accented character.

## Usage

Most of the command line documentation can be fetch through `primetools help`
command.

### Dumping crates/playlists to files

Let's you want to dump crates saved on a external disk (export) located on the P
drive which have 'fancy' in their name/path.

```bash
primetools dump crates -sp "p:" -o crates.yaml -n "*fancy*"
```

Dump Help:

```txt
NAME:
   primetools dump - the swiss knife of Denon's Engine PRIME

USAGE:
   primetools dump [command options] [arguments...]

DESCRIPTION:
   dump data about a library

OPTIONS:
   --source value, -s value         (default: ITunes)
   --source-path value, --sp value
   --output value, -o value         (default: "-")
   --format value, -f value         (default: Auto)
   --name value, -n value
```

### Fixing missing file

Let's say you moved files around and want to fix those files. This will search
in the specified folder for a file that matches the same meta data. If the meta
data has changed, it won't works. Also, if more than one match is found, the
program will ask which one you want to use as a fix.

```bash
primetools fix missing -s prime -p M:\\super\\folder\\to\\search
```

```txt
USAGE:
    primetools fix [command options] [arguments...]

DESCRIPTION:
    try to fix problem database [Duplicate, Missing]

OPTIONS:
    --source value, -s value         (default: ITunes)
    --source-path value, --sp value
    --yes, -y                        Do not prompt for write confirmation (default: false)
    --search-path value, -p value    path to search for music file
    --dryrun, --ro                   (default: false)
```

### Syncing

If you want to sync the `added` date from ITunes to PRIME.

```bash
primetools sync added -s itunes -t prime
```

```txt
USAGE:
   primetools sync [command options] [arguments...]

DESCRIPTION:
   sync assets from a source to a destination [Ratings, Added, Modified, PlayCount]

OPTIONS:
   --source value, -s value         (default: ITunes)
   --source-path value, --sp value
   --target value, -t value         (default: PRIME)
   --target-path value, --tp value
   --dryrun, --ro                   (default: false)
```

### Importing crates / playlist

You can import crates/playlist from . Note that if a list already exists, its
content _will be overriden_. The tools doesn't support merging.

```bash
primetools import -s crates.yaml
```

```txt
USAGE:
   primetools import [command options] [arguments...]

DESCRIPTION:
   import playlist/crates

OPTIONS:
   --target value, -t value         (default: PRIME)
   --target-path value, --tp value
   --source value, -s value         file to use as source
   --name value, -n value           Names of crate/playlist to import, can be glob (*something*), 
                                    if none is given, will import all object in dump file.
   --ignore-not-found               Ignore track which aren't found in target, otherwise the 
                                    operation will fail. (default: false)
   --dryrun, --ro                   (default: false)
```

## Ref

- [Engine Library Format](https://github.com/mixxxdj/mixxx/wiki/engine_library_format)
