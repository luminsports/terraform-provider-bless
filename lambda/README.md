Update dependencies:
  
```pip3.10 install \
--platform manylinux2014_x86_64 \
--target=. \
--implementation cp \
--python-version 3.10 \
--only-binary=:all: --upgrade \
-r requirements.txt```
