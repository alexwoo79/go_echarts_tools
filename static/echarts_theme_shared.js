(function () {
  var PROFILES = {
    // Theme builder JSON: public/themes/default.json
    "default": {
      palette: ["#c23531", "#2f4554", "#61a0a8", "#d48265", "#91c7ae", "#749f83", "#ca8622", "#bda29a", "#6e7074", "#546570", "#c4ccd3"],
      backgroundColor: "rgba(0,0,0,0)",
      titleColor: "#333333",
      subtitleColor: "#aaa",
      textColor: "#333",
      axisLineColor: "#333",
      axisLabelColor: "#333",
      splitLineColor: "#ccc",
      toolboxColor: "#999999",
      toolboxEmphasisColor: "#666666",
      tooltipAxisColor: "#cccccc",
      isDark: false
    },
    // Theme builder JSON: public/themes/v5.json
    "v5": {
      palette: ["#5470c6", "#91cc75", "#fac858", "#ee6666", "#73c0de", "#3ba272", "#fc8452", "#9a60b4", "#ea7ccc"],
      backgroundColor: "rgba(0,0,0,0)",
      titleColor: "#464646",
      subtitleColor: "#6E7079",
      textColor: "#333",
      axisLineColor: "#6E7079",
      axisLabelColor: "#6E7079",
      splitLineColor: "#E0E6F1",
      toolboxColor: "#999",
      toolboxEmphasisColor: "#666",
      tooltipAxisColor: "#ccc",
      isDark: false
    },
    // Theme builder JSON: public/themes/dark.json
    "dark": {
      palette: ["#dd6b66", "#759aa0", "#e69d87", "#8dc1a9", "#ea7e53", "#eedd78", "#73a373", "#73b9bc", "#7289ab", "#91ca8c", "#f49f42"],
      backgroundColor: "rgba(51,51,51,1)",
      titleColor: "#eeeeee",
      subtitleColor: "#aaa",
      textColor: "#eeeeee",
      axisLineColor: "#eeeeee",
      axisLabelColor: "#eeeeee",
      splitLineColor: "#aaaaaa",
      toolboxColor: "#999",
      toolboxEmphasisColor: "#666",
      tooltipAxisColor: "#eeeeee",
      isDark: true
    },
    // Theme builder JSON: public/themes/vintage.json
    "vintage": {
      palette: ["#d87c7c", "#919e8b", "#d7ab82", "#6e7074", "#61a0a8", "#efa18d", "#787464", "#cc7e63", "#724e58", "#4b565b"],
      backgroundColor: "#fef8ef"
    },
    // Theme builder JSON: public/themes/westeros.json
    "westeros": {
      palette: ["#516b91", "#59c4e6", "#edafda", "#93b7e3", "#a5e7f0", "#cbb0e3"],
      backgroundColor: "transparent"
    },
    // Theme builder JSON: public/themes/essos.json
    "essos": {
      palette: ["#893448", "#d95850", "#eb8146", "#ffb248", "#f2d643", "#ebdba4"],
      backgroundColor: "rgba(242,234,191,0.15)"
    },
    // Theme builder JSON: public/themes/wonderland.json
    "wonderland": {
      palette: ["#4ea397", "#22c3aa", "#7bd9a5", "#d0648a", "#f58db2", "#f2b3c9"],
      backgroundColor: "transparent"
    },
    // Theme builder JSON: public/themes/walden.json
    "walden": {
      palette: ["#3fb1e3", "#6be6c1", "#626c91", "#a0a7e6", "#c4ebad", "#96dee8"],
      backgroundColor: "rgba(252,252,252,0)"
    },
    // Theme builder JSON: public/themes/chalk.json
    "chalk": {
      palette: ["#fc97af", "#87f7cf", "#f7f494", "#72ccff", "#f7c5a0", "#d4a4eb", "#d2f5a6", "#76f2f2"],
      backgroundColor: "#293441",
      isDark: true
    },
    // Theme builder JSON: public/themes/infographic.json
    "infographic": {
      palette: ["#c1232b", "#27727b", "#fcce10", "#e87c25", "#b5c334", "#fe8463", "#9bca63", "#fad860", "#f3a43b", "#60c0dd", "#d7504b", "#c6e579", "#f4e001", "#f0805a", "#26c0c0"],
      backgroundColor: "rgba(0,0,0,0)",
      titleColor: "#27727b",
      toolboxColor: "#c1232b",
      toolboxEmphasisColor: "#e87c25"
    },
    // Theme builder JSON: public/themes/macarons.json
    "macarons": {
      palette: ["#2ec7c9", "#b6a2de", "#5ab1ef", "#ffb980", "#d87a80", "#8d98b3", "#e5cf0d", "#97b552", "#95706d", "#dc69aa", "#07a2a4", "#9a7fd1", "#588dd5", "#f5994e", "#c05050", "#59678c", "#c9ab00", "#7eb00a", "#6f5553", "#c14089"],
      backgroundColor: "rgba(0,0,0,0)",
      titleColor: "#008acd",
      toolboxColor: "#2ec7c9",
      toolboxEmphasisColor: "#18a4a6"
    },
    // Theme builder JSON: public/themes/roma.json
    "roma": {
      palette: ["#e01f54", "#001852", "#f5e8c8", "#b8d2c7", "#c6b38e", "#a4d8c2", "#f3d999", "#d3758f", "#dcc392", "#2e4783", "#82b6e9", "#ff6347", "#a092f1", "#0a915d", "#eaf889", "#6699FF", "#ff6666", "#3cb371", "#d5b158", "#38b6b6"],
      backgroundColor: "rgba(0,0,0,0)"
    },
    // Theme builder JSON: public/themes/shine.json
    "shine": {
      palette: ["#c12e34", "#e6b600", "#0098d9", "#2b821d", "#005eaa", "#339ca8", "#cda819", "#32a487"],
      backgroundColor: "transparent"
    },
    // Theme builder JSON: public/themes/purple-passion.json
    "purple-passion": {
      palette: ["#8a7ca8", "#e098c7", "#8fd3e8", "#71669e", "#cc70af", "#7cb4cc"],
      backgroundColor: "rgba(91,92,110,1)",
      isDark: true
    },
    // Theme builder JSON: public/themes/halloween.json
    "halloween": {
      palette: ["#ff715e", "#ffaf51", "#ffee51", "#8c6ac4", "#715c87"],
      backgroundColor: "rgba(64,64,64,0.5)",
      titleColor: "#ffaf51",
      subtitleColor: "#eeeeee",
      axisLineColor: "#666666",
      axisLabelColor: "#999999",
      splitLineColor: "#555555",
      toolboxColor: "#999999",
      toolboxEmphasisColor: "#666666",
      tooltipAxisColor: "#cccccc",
      isDark: true
    }
  };

  var CDN_REGISTERED_THEMES = {
    dark: true,
    vintage: true,
    macarons: true,
    shine: true,
    roma: true,
    infographic: true
  };

  function normalizeThemeName(name) {
    var n = String(name || "default").toLowerCase();
    if (n === "v5") return "default";
    return PROFILES[n] ? n : "default";
  }

  function mergeProfile(profile) {
    var def = PROFILES.default;
    return {
      name: profile.name || "default",
      palette: profile.palette || def.palette,
      backgroundColor: profile.backgroundColor != null ? profile.backgroundColor : def.backgroundColor,
      titleColor: profile.titleColor || def.titleColor,
      subtitleColor: profile.subtitleColor || def.subtitleColor,
      textColor: profile.textColor || def.textColor,
      axisLineColor: profile.axisLineColor || def.axisLineColor,
      axisLabelColor: profile.axisLabelColor || def.axisLabelColor,
      splitLineColor: profile.splitLineColor || def.splitLineColor,
      toolboxColor: profile.toolboxColor || def.toolboxColor,
      toolboxEmphasisColor: profile.toolboxEmphasisColor || def.toolboxEmphasisColor,
      tooltipAxisColor: profile.tooltipAxisColor || def.tooltipAxisColor,
      isDark: !!profile.isDark
    };
  }

  function getThemeProfile(name) {
    var normalized = normalizeThemeName(name);
    var src = PROFILES[normalized] || PROFILES.default;
    var withName = Object.assign({ name: normalized }, src);
    return mergeProfile(withName);
  }

  function getThemeNames() {
    return Object.keys(PROFILES).filter(function (name) { return name !== "v5"; });
  }

  function getEchartsThemeName(name) {
    var normalized = normalizeThemeName(name);
    if (normalized === "default") return null;
    return CDN_REGISTERED_THEMES[normalized] ? normalized : null;
  }

  window.EChartsThemeShared = {
    getThemeProfile: getThemeProfile,
    getThemeNames: getThemeNames,
    getEchartsThemeName: getEchartsThemeName,
    normalizeThemeName: normalizeThemeName
  };
})();
