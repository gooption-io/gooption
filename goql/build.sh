g++ \
    -std=c++14 -stdlib=libc++ \
    -O2 -Wall \
    -I /usr/local/include/boost \
    -I ../proto/cpp \
    -L/usr/local/lib \
    -lQuantLib \
    -lboost_program_options \
    -lgrpc++_reflection \
    -lprotobuf -lgrpc++ -lgrpc \
    ../proto/cpp/*.cc goql.cc \
    -o goql