## Running locally

You will need to use environment variables in a `.env.local` file to run the project. The best way is to get them directly from the Vercel project with `vercel env pull`

Run
```
npx vercel env pull .env.local --yes --environment=development --token=[VERCEL_TOKEN]
```
Replace `VERCEL_TOKEN` with the Vercel token for the account with access to the project.

Once you have the `.env.local` file, you can start the dev sever by running
```
npm run build
npm next dev --turbo -p [PORT]
```
Replace `PORT` with the port you'd like the server to run on, if you don't include the `-p` flag it will default to port 3000.

The app  should now be running on [localhost:3000](http://localhost:3000/).
