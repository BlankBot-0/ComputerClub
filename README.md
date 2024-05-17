Инструкции к запуску (в директории проекта):
1. $ docker build -t test-cont .
2. $ docker run --rm -v /your/path/to/tests/:root/files test-cont your_test_file.txt