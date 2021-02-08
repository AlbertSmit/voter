const sveltePreprocess = require("svelte-preprocess");
const typescript = require("@rollup/plugin-typescript");
const preprocess = sveltePreprocess({
  postcss: {
    plugins: [require("tailwindcss"), require("autoprefixer")],
  },
});

const ts = typescript({ sourceMap: !process.env.production });

module.exports = {
  preprocess,
  ts,
};
