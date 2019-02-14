g++ \
    -std=c++14 -stdlib=libc++ \
    -O2 -Wall \
    -I /usr/local/include/boost \
    -I /usr/local/include/spdlog \
    -I . \
    -L/usr/local/lib \
    -lboost_program_options \
    -lgrpc++_reflection \
    -lprotobuf -lgrpc++ -lgrpc \
    *.cc \
    -o goql \
    /usr/local/lib/libQuantLib.a
