alter table if exists events
    drop column is_notified,
    drop column when_to_notify,
    add column when_to_notify varchar(256),
    alter column when_to_notify drop not null,
    alter column description drop not null,
    alter column description drop default
;



