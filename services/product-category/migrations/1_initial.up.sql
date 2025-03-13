create or replace function update_updated_at_column()
    returns trigger as $$
begin
    new.modified = now();
    return new;
end;
$$ language plpgsql;


