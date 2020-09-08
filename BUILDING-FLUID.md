sudo apt install libasound2-dev fsutils
cmake -Denable-jack=off -Denable-sdl2=off -DLIB_SUFFIX="" .
sudo usermod -a -G floppy coral
sudo usermod -a -G disk coral
