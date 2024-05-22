SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: client_scopes; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.client_scopes (
    id integer NOT NULL,
    client_id integer,
    resource character varying(150) NOT NULL,
    scopes jsonb DEFAULT '[]'::jsonb,
    disabled boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,
    disabled_at timestamp with time zone
);


--
-- Name: client_scopes_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.client_scopes_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: client_scopes_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.client_scopes_id_seq OWNED BY public.client_scopes.id;


--
-- Name: clients; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.clients (
    id integer NOT NULL,
    name character varying(25) NOT NULL,
    secret character varying(128) NOT NULL,
    expired_at timestamp with time zone,
    disabled boolean DEFAULT false,
    created_at timestamp with time zone DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone,
    disabled_at timestamp with time zone
);


--
-- Name: clients_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.clients_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: clients_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.clients_id_seq OWNED BY public.clients.id;


--
-- Name: schema_migrations; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.schema_migrations (
    version character varying(128) NOT NULL
);


--
-- Name: client_scopes id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_scopes ALTER COLUMN id SET DEFAULT nextval('public.client_scopes_id_seq'::regclass);


--
-- Name: clients id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.clients ALTER COLUMN id SET DEFAULT nextval('public.clients_id_seq'::regclass);


--
-- Name: client_scopes client_scopes_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_scopes
    ADD CONSTRAINT client_scopes_pkey PRIMARY KEY (id);


--
-- Name: clients clients_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.clients
    ADD CONSTRAINT clients_pkey PRIMARY KEY (id);


--
-- Name: schema_migrations schema_migrations_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.schema_migrations
    ADD CONSTRAINT schema_migrations_pkey PRIMARY KEY (version);


--
-- Name: client_scopes_client_id_resource_uidx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX client_scopes_client_id_resource_uidx ON public.client_scopes USING btree (client_id, resource);


--
-- Name: clients_name_uidx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX clients_name_uidx ON public.clients USING btree (name);


--
-- Name: clients_secret_uidx; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX clients_secret_uidx ON public.clients USING btree (secret);


--
-- Name: client_scopes fk_client_scopes_client_id; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.client_scopes
    ADD CONSTRAINT fk_client_scopes_client_id FOREIGN KEY (client_id) REFERENCES public.clients(id);


--
-- PostgreSQL database dump complete
--


--
-- Dbmate schema migrations
--

INSERT INTO public.schema_migrations (version) VALUES
    ('20240515015140');
