# Alerts

Alerts are modelled with the following fields:

- subject:
    - will be saved inside the database so the user is able to know which alert was sent/not sent
- status:
    - pending/sent/failure
- sent_at:
    - timestamp which visualizes when the mail was sent
- user_id:
    - connection between email and user
- created_at
- updated_at
