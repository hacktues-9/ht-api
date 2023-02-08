create view teams
            (id, name, logo, project_id, approved)
as
SELECT team.id,
       team.name,
       team.logo,
       team.project_id,
       team.approved
FROM team
WHERE team.deleted_at IS NULL
ORDER BY team.created_at;

alter table teams
    owner to doadmin;
