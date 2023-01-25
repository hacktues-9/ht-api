create view teams
            (tid, name, logo, pid, uid, first_name, last_name, email, profile_picture, role, grade, class, discord, github, approved)
as
SELECT team.id         AS tid,
       team.name,
       team.logo,
       team.project_id AS pid,
       users.id        AS uid,
       users.first_name,
       users.last_name,
       users.email,
       socials.profile_picture,
       role.name       AS role,
       info.grade,
       class.name      AS class,
       concat(discord.username, '#', discord.discriminator) AS discord,
       github.login AS github,
       team.approved
FROM users
         LEFT JOIN team ON users.team_id = team.id
         JOIN info ON users.info_id = info.id
         JOIN socials ON info.socials_id = socials.id
         LEFT JOIN class ON info.class_id = class.id
         LEFT JOIN discord ON socials.discord_id = discord.id
         LEFT JOIN github ON socials.github_id = github.id
         LEFT JOIN role ON users.role_id = role.id
WHERE team.id IS NOT NULL;

alter table teams
    owner to postgres;
