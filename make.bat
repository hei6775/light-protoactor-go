@set PACKAGE_PATH=gitee.com\lwj8507\light-protoactor-go
@set WORK_DIR=%GOPATH%\src\%PACKAGE_PATH%
@set VENDOR_DIR=%WORK_DIR%\vendor


@IF "%1" == "glide-up" call :glide-up & goto :exit

@IF "%1" == "zip-vendor" call :zip-vendor & goto :exit

@IF "%1" == "unzip-vendor" call :unzip-vendor & goto :exit

@echo unsupported operate [%1]

@goto :exit


:glide-up
@echo glide up begin
glide up --strip-vendor
@echo glide up end
@goto :exit


:zip-vendor
@echo zip [vendor] begin
del "%VENDOR_DIR%.zip"
del "%VENDOR_DIR%\gitee.com\lwj8507\nggs\vendor.zip"
7z a -tzip "%VENDOR_DIR%.zip" "%VENDOR_DIR%"
@echo zip [vendor] end
@goto :exit


:unzip-vendor
@echo unzip [vendor] begin
rmdir /q /s "%WORK_DIR%\vendor"
7z x "%WORK_DIR%\vendor.zip" -y -aos -o"%WORK_DIR%"
@echo unzip [vendor] end
@goto :exit


:exit
