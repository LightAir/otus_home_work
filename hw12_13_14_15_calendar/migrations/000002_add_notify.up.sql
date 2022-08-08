alter table if exists events
    add column is_notified bool default false not null,
    drop column when_to_notify,
    add column when_to_notify timestamptz not null,
    alter column description set not null,
    alter column description set default ''
;