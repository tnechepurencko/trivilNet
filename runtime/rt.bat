clang rt_api.c -ldbghelp --shared -fno-omit-frame-pointer -o libwelrt.dll
rem linux: *files* -fPIC --shared -fno-omit-frame-pointer -ldl -lm -o libwelrt.so