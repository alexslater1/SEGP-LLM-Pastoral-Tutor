CREATE TABLE IF NOT EXISTS "RagDocument" (
	"id" uuid PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL,
	"name" text NOT NULL,
	"uploadedAt" timestamp DEFAULT now() NOT NULL,
	"size" text NOT NULL,
	"type" varchar NOT NULL,
	"status" varchar DEFAULT 'processing' NOT NULL,
	"userId" uuid NOT NULL,
	"backendSourceId" integer
);
--> statement-breakpoint
DO $$ BEGIN
 ALTER TABLE "RagDocument" ADD CONSTRAINT "RagDocument_userId_User_id_fk" FOREIGN KEY ("userId") REFERENCES "public"."User"("id") ON DELETE no action ON UPDATE no action;
EXCEPTION
 WHEN duplicate_object THEN null;
END $$;
