<script lang="ts">
  import { onMount } from "svelte";
  import * as PlatformService from "../../../bindings/TcNo-Acc-Switcher/internal/platform/platformservice.js";
  import { openExternalUrl } from "../../lib/openExternalUrl";
  import { t } from "../../stores/i18n";
  import { checkForUpdatesManually } from "../../lib/checkForUpdates";

  import BrqLogo from "../BrqLogo.svelte";

  const UPSTREAM_REPO_URL = "https://github.com/TCNOco/TcNo-Acc-Switcher";
  const OPEN_SOURCE_LICENSES_URL = "https://github.com/TCNOco/TcNo-Acc-Switcher/blob/master/OPEN_SOURCE_LICENSES.txt";

  let currentVersion = "0.0.0";

  onMount(() => {
    void PlatformService.GetAppVersion()
      .then((v) => {
        currentVersion = v || "0.0.0";
      })
      .catch(() => {
        currentVersion = "0.0.0";
      });
  });

  function openExternal(url: string, e: MouseEvent): void {
    e.preventDefault();
    void openExternalUrl(url);
  }

  async function onVersionClick(): Promise<void> {
    await checkForUpdatesManually();
  }
</script>

<div class="about-modal">
  <div class="imgDiv">
    <BrqLogo variant="mark" class="about-mark" />
  </div>
  <div class="rightContent">
    <h2>BRQ Account Switcher</h2>
    <p>
      Based on <strong>TcNo Account Switcher</strong> by TroubleChute, used under the GPL-3.0 licence.
      BRQ is an independent edition and is not affiliated with TcNo or TroubleChute.
    </p>
    <div class="linksList">
      <a href={UPSTREAM_REPO_URL} on:click={(e) => openExternal(UPSTREAM_REPO_URL, e)}>
        <svg viewBox="0 0 24 24" aria-hidden="true" class="modalIcoGitHub">
          <use href="img/icons/ico_github.svg#icoGitHub"></use>
        </svg>
        Original project (TcNo Account Switcher)
      </a>
      <a
        href="https://github.com/TCNOCo/TcNo-Acc-Switcher/blob/master/DISCLAIMER.md"
        on:click={(e) =>
          openExternal("https://github.com/TCNOCo/TcNo-Acc-Switcher/blob/master/DISCLAIMER.md", e)}
      >
        <svg viewBox="0 0 2084 2084" aria-hidden="true" class="modalIcoDoc">
          <use href="img/icons/ico_doc.svg#icoDoc"></use>
        </svg>
        {$t("Modal_Info_Disclaimer")}
      </a>
      <a href={OPEN_SOURCE_LICENSES_URL} on:click={(e) => openExternal(OPEN_SOURCE_LICENSES_URL, e)}>
        <svg viewBox="0 0 2084 2084" aria-hidden="true" class="modalIcoDoc">
          <use href="img/icons/ico_doc.svg#icoDoc"></use>
        </svg>
        Open source licenses
      </a>
    </div>
  </div>
</div>
<div class="versionIdentifier">
  <span>{$t("Modal_Info_Version")}: <button type="button" class="fancyLink version-link" on:click={() => void onVersionClick()}>{currentVersion}</button></span>
</div>

<style lang="scss">
  .about-modal {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    gap: 1rem;
    align-items: flex-start;
  }

  .imgDiv {
    flex: 0 0 auto;
  }

  .imgDiv :global(.about-mark) {
    width: 100px;
    height: 100px;
  }

  .rightContent {
    flex: 1 1 200px;
    min-width: 0;
  }

  .rightContent h2 {
    margin: 0 0 0.5rem;
    font-size: 1.25rem;
  }

  .rightContent p {
    margin: 0 0 0.75rem;
  }

  .linksList {
    display: flex;
    flex-direction: column;
    gap: 0.35rem;
  }

  .linksList svg {
    flex-shrink: 0;
    margin-right: 0.25rem;
    width: 1.5rem;
    height: 1.5rem;

    use {
      fill: var(--whiteSecondary);
    }
  }

  .versionIdentifier {
    margin-top: 0.75rem;
    padding-top: 0.5rem;
    border-top: 1px solid var(--border-bar-bg);
    font-size: 0.85rem;
    opacity: 0.85;
  }

  .version-link {
    text-decoration: underline;
  }
</style>
