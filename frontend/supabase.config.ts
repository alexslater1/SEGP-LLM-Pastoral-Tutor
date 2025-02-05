const config = {
    appName: "Pastoral Tutor",
    appDescription:
      "Pastoral Tutor AI",
    domainName:
      process.env.NODE_ENV === "development"
        ? "http://localhost:3000"
        : "http://localhost:3000",
  };
  
  export default config;