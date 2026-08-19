-- Applied after every migration so River's own tables are covered too.
grant usage on schema public to getnotified_app;
grant all on all tables in schema public to getnotified_app;
grant all on all sequences in schema public to getnotified_app;
grant execute on all functions in schema public to getnotified_app;
