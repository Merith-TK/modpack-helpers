# how to use
This repo contains a hanful of tools

# how to install

Make sure you have golang installed (https://golang.org/doc/install),
and then run
`go install github.com/Merith-TK/modpack-helpers/cmd/<tool>@main`


## Mod Deduplicator
`modjar-deduplicator`
- Usage
  - open terminal where the mods folder is
  - this would be the `.minecraft` folder
  `modjar-deduplicator` and it will automatically rename jars
  

## Packwiz Wrapper
`pw`
- Usage

  - `pw version`
    - prints pw version
  - `pw help`
    - prints help
  - `pw import <import.txt>`
    - imports mods from links in import.txt
    - currently supports importing from curseforge, modrinth, and URL files (URL files have formatting)
```md
https://modrinth.com/mod/AANobbMI
https://www.curseforge.com/minecraft/mc-mods/lithium

This will download mardown-test to example.md in the pack
/example.md
https://cdn.merith.xyz/markdown-test.md
```
  - `pw modlist [raw] [versions]`
    - Generates a modlist file
    - arguments are keyword matches, so you can do `pw modlist raw versions` or `pw modlist versions raw`
    - You can use `raw` to generate without markdown formatting, (can be used for import)
    - You can use `versions` to generate with versions specified in modlist
    - This is automatically sorted (where possible) for client, shared, and server mods,
      - Shared mods are mods that are required on both client and server for full functionality
      - this is detected from the `side` field in the mod's `pack.toml`

  - `pw reinstall`
    - reinstalls all meta-files, including URL files
  - `pw batch (packwiz args)`
    - same as old pw -b, runs packwiz command in all subfolders of a folder
    - uses `./` by default but with `pw -d subdir batch` it will use `subdir` folder instead
    - Example: `pw batch refresh` will refresh all subpacks

- `pw` supports a few flags, which can be used with any subcommand
  - `pw -h` will print help
  - `pw -r` will run `packwiz refresh` after operations
  - running `pw -r import` will automatically refresh after its done importing
  - `pw -y` will autoconfirm (not full implemented into all subcommands)
  - `pw -c` is depreciated
  - Originally used for importing only clientside mods, but does nothing now
  - `pw -d <PackDir>` will set the pack directory to `<PackDir>`
  - In batchmode this points to where folders containg pack.toml's are rather than the current folder




