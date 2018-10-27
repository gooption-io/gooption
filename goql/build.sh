g++ \
    -std=c++14 -stdlib=libc++ \
    -O2 -fno-exceptions -Wall \
    -I /usr/local/include/boost \
    -I ../proto/cpp \
    -L/usr/local/lib \
    -lQuantLib \
    -lgrpc++_reflection \
    -lprotobuf -lgrpc++ -lgrpc \
    ../proto/cpp/*.cc goql_server.cc \
    -o goql