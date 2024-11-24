# Wa-Tor-Simulation

A GO simulation of the Wa-Tor world.

### Before running or compiling the program make sure that [Fyne](https://docs.fyne.io/started/) is set up correctly.

To run the simulation type in the following to a terminal in the **source** directory:

`go run .`

To build and run executable:

`go build . && ./source`

To view documentation do one of the following steps:

1. Open up the **documentation** directory in a HTML server such as VS Code Live Server.

2. Install the [Golds](https://github.com/go101/golds) package using:

   `go install go101.org/golds@latest`

   Then change to the **source** directory and run:

   `golds .`

   If you get an error saying `golds` cannot be found it means that the Go binary installation path is not in the `PATH` environment variable. In that case either add it or run from where it is installed, usually:

   `~/go/bin/golds .`

Following one of the above 2 steps a web browser should then open on the `index.html` page. Sort by **dependancy distance** and select the first one to read the documentation of the code i.e. **github.com/NotYetTerminal/Wa-Tor-Simulation/source**

# LICENCE

All files not declared with a license are under the following license:  
Concurrent-Development © 2024 by Gábor Major is licensed under CC BY-NC-SA 4.0. To view a copy of this license,
visit https://creativecommons.org/licenses/by-nc-sa/4.0/
