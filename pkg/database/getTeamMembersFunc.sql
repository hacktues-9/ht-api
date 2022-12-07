-- Active: 1663917833212@@db-postgresql-ht-9-test-do-user-12488008-0.b.db.ondigitalocean.com@25060@defaultdb@public
create function getTeamMembers(teamID BIGINT)
returns table (userID BIGINT)
language plpgsql
AS
$$
declare
    teamMembers table (userID BIGINT);
begin
    select id
    into teamMembers
    from user
    where team_id = teamID;

    return teamMembers;
end;
$$;