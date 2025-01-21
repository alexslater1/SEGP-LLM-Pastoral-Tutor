## Running locally

### Installing Packages
You need to install the package dependencies first using `pnpm` (you can install `pnpm` by following the instructions [here](https://pnpm.io/installation).

Once you've installed it, run
```
pnpm install
```

### Getting Environment Variables
You will need to use environment variables in a `.env.local` file to run the project. The best way is to get them directly from the Vercel project with `vercel env pull`

Run
```
export VERCEL_ORG_ID=team_zI0C0C3bsu0Tu6wPEVRr4AnP
export VERCEL_PROJECT_ID=team_zI0C0C3bsu0Tu6wPEVRr4AnP
npx vercel env pull .env.local --yes --environment=development --token=[VERCEL_TOKEN]
```
Replace `VERCEL_TOKEN` with the Vercel token for the account with access to the project.

Once you have the `.env.local` file, you can start the dev sever by running
```
npm run build
npx next dev --turbo -p [PORT]
```
Replace `PORT` with the port you'd like the server to run on, if you don't include the `-p` flag it will default to port 3000.

The app  should now be running on [localhost:3000](http://localhost:3000/).
