Инструкции к запуску (в директории проекта):
```sh
$ docker build -t computer-club .
$ docker run --rm -v /your/path/to/tests/:root/files test-cont your_test_file.txt
```