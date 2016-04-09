恢复备份
mongorestore.exe --db GoodShare C:\Users\payallmoney\Downloads\GoodShare\GoodShare


备份数据库
mongodump --db GoodShare --out ~/mongodbbackup
zip -r GoodShare.zip ~/mongodbbackup/GoodShare

