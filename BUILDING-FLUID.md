sudo apt install libasound2-dev fdutils libglib2.0-dev libsndfile-dev
cmake -Denable-jack=off -Denable-sdl2=off -DLIB_SUFFIX="" .
sudo usermod -a -G floppy,disk coral
