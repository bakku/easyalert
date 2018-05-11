# Emails

Emails are modelled with the following fields:

- name:
    - will be saved inside the database so the user is able to know which mail was sent/not sent
- status:
    - pending/sent/failure
- sent_at:
    - timestamp which visualizes when the mail was sent
- user_id:
    - connection between email and user
- created_at
- updated_at