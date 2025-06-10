"# igorm" 
dbx row_number example:
select row_number() as stt * from table_name order by column_name
-> select row_number() over (order by column_name) as stt * from table_name;
run swag init -g cmd/api/main.go --parseDependency