I want to change the nature of how binary builds are created in this project (inspired_trek73) to match the same manner in which they are handled in the Inspired_Trek73_Multiplayer project.

Make a plan for the following. Do not change any code until the plan is approved.

Specifically, there will be a single build script in a subdirectory called "scripts" that runs under windows powershell that when executed creates all of the target binaries and puts them in a directory called "dist." It should target the same operating systems and architectures as currently targeted by the Inspired_Treki73_Multiplayer project.

To be clear, both the scripts directory and the dist directory should be immediatley below the top level directory for the project (namely right below Inspired_Trek73).