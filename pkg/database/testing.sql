-- PostgreSQL 

-- view with data from tables - User and Info

CREATE OR REPLACE VIEW user_info AS
SELECT u.id, u.name, u.email, i.phone, i.address
FROM users u
INNER JOIN info i ON u.id = i.user_id;
