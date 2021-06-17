@echo off
mkdir "C:\Program Files (x86)\SmbUploader"
mkdir "C:\Program Files (x86)\SmbUploader\img"

copy uploader.exe "C:\Program Files (x86)\SmbUploader"
copy config.ini "C:\Program Files (x86)\SmbUploader"
copy data.db.example "C:\Program Files (x86)\SmbUploader\data.db"
copy setup.reg "C:\Program Files (x86)\SmbUploader"
xcopy img "C:\Program Files (x86)\SmbUploader\img"

echo "Run setup.reg"
regedit /s "C:\Program Files (x86)\SmbUploader\setup.reg"
exit