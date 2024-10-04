# NESEmuZ

# NES Emulator in Go

This project is a NES (Nintendo Entertainment System) emulator written in Go. The emulator aims to faithfully reproduce the behavior of the original NES hardware, allowing you to play classic NES games on your computer.

## Features

- Accurate CPU emulation
- PPU (Picture Processing Unit) emulation for rendering graphics
- APU (Audio Processing Unit) emulation for sound
- Support for various mappers
- Save and load game states
- Controller input handling

## Getting Started

### Prerequisites

- Go 1.16 or later
- A terminal or command prompt

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/lizi-cl/NESEmuZ.git
    cd NESEmuZ
    ```

2. Build the project:
    ```sh
    go build -o NESEmuZ
    ```

3. Run the emulator:
    ```sh
    ./NESEmuZ path/to/your/rom.nes
    ```

### Usage

To run the emulator, you need to provide the path to a NES ROM file. You can specify the ROM file as a command-line argument:

```sh
./NESEmuZ path/to/your/rom.nes