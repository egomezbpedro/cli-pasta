# CLI-Pasta a fuzzy clipboard finder

Fuzzy find your clips


# Installation

Install the background deamon that listen for clipboard updates
```sh
    go install github.com/egomezbpedro/cli-pasta/pasta-deamon@latest \
    pasta-deamon install \
    pasta-deamon satus
    
```

Install the cli interface
```sh
    go install github.com/egomezbpedro/cli-pasta@latest
```

# Limitations

- Service deamon only works for OSx systems [WIP].

# Usage

![ezgif com-video-to-gif (1)](https://github.com/egomezbpedro/cli-pasta/assets/57415533/d57b19bc-1890-4ef7-9ad1-3358e317bf53)

**Neovim remap**
``` lua
    vim.keymap.set("n", "<leader>cp", "<Cmd>:silent !tmux split-window -h cli-pasta<CR>")
```
