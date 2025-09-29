
## Prepare before development

1. Create file .env or run this command in root project

```bash
  touch .env 
```

2. Set variable like in this example or you can follow in .env.example

```bash
GO_ENV=development
PORT=8000
DB_URL=postgresql://postgres:$USERNAME@HOST:PORT/DATABASE
```

3. Open your terminal and run command Air package
```bash
air
```

4. Your terminal will see like this 
```bash
2025/09/29 13:59:00 ##### Begin load ENV #####
2025/09/29 13:59:00 Service run in mode development
2025/09/29 13:59:00 Successfully connected to database pgx ✅!

   ____    __
  / __/___/ /  ___
 / _// __/ _ \/ _ \
/___/\__/_//_/\___/ v4.13.4
High performance, minimalist Go web framework
https://echo.labstack.com
____________________________________O/_______
                                    O\
⇨ http server started on [::]:8000
```

5. Use command for generate privateKey, public key to sign and verify JWT token
```bash
// keySize 1024
openssl genrsa -out app.rsa ${keySize}
openssl rsa -in app.rsa -pubout > app.rsa.pub
```

## Thank you distributors and document at
- air: https://github.com/air-verse/air
- golang: https://go.dev/
- golang echo framework: https://echo.labstack.com
- Squirrel - fluent SQL generator for Go: https://github.com/Masterminds/squirrel
- pgx - PostgreSQL Driver and Toolkit: https://github.com/jackc/pgx

## If you have Questions you can contact me
- Githubs: https://github.com/Binh-2060
- Website: https://binh.phouservice.com
- Email: binhxayxana@gmail.com



    