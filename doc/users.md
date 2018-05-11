# Users

Users are modelled with the following fields:

- email: 
    - can be used for authentication in combination with password via HTTP Basic
    - will be used as target address when sending an email
- password:
    - stored with bcrypt
    - can be used for authentication in combination with email via HTTP Basic
- token:
    - is generated on signup
    - can be used for authentication if user does not want to expose email and password
- admin:
    - has access to all data
- created_at
- updated_at