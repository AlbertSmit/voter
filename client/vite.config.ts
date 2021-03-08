import svelte from "rollup-plugin-svelte";
import sveltePreprocess from "svelte-preprocess";
import autoprefixer from "autoprefixer";

const sourceMapsInProduction = false;

const production = process.env.NODE_ENV === "production";
module.exports = {
  root: "./src",
  plugins: [
    svelte({
      emitCss: production,
      preprocess: sveltePreprocess(),
      compilerOptions: {
        dev: !production,
      },

      // @ts-ignore This is temporary until the type definitions are fixed!
      hot: !production,
    }),
  ],
  server: {
    host: "localhost",
    port: 5000,
  },
  build: {
    outDir: "../dist",
    sourcemap: sourceMapsInProduction,
  },
};
