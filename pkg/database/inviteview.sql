create function inviteview(uid bigint)
    returns TABLE( team_id bigint, team_name text, team_logo text)
    language plpgsql
as
$$
BEGIN
    RETURN QUERY
        select team.id, team.name, team.logo
        from invite
                 join team on team.id = invite.team_id
        where invite.user_id = uid;
END;
$$;

alter function inviteview(bigint) owner to doadmin;

