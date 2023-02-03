create view teams
            (id, name, logo, project_id, approved)
as
SELECT team.id,
       team.name,
       team.logo,
       team.project_id,
       team.approved
FROM team
ORDER BY team.created_at DESC;

alter table teams
    owner to doadmin;
