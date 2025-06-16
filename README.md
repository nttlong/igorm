"# igorm" 
dbx row_number example:
select row_number() as stt * from table_name order by column_name
-> select row_number() over (order by column_name) as stt * from table_name;
run swag init -g cmd/api/main.go --parseDependency
npm create vite@latest x-unvs -- --template react
cd x-unvs
npm install
npm run dev
npm install -D tailwindcss@3 postcss@8 autoprefixer@10