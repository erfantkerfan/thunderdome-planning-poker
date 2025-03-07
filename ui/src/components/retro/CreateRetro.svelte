<script lang="ts">
  import { onMount } from 'svelte';

  import SolidButton from '../SolidButton.svelte';
  import DownCarrotIcon from '../icons/ChevronDown.svelte';
  import { warrior as user } from '../../stores';
  import { AppConfig, appRoutes } from '../../config';
  import LL from '../../i18n/i18n-svelte';

  export let xfetch;
  export let notifications;
  export let eventTag;
  export let router;
  export let apiPrefix = '/api';

  let retroName = '';
  let joinCode = '';
  let facilitatorCode = '';
  let maxVotes = 3;
  let brainstormVisibility = 'visible';
  let teams = [];
  let selectedTeam = '';

  const brainstormVisibilityOptions = [
    {
      label: $LL.brainstormVisibilityLabelVisible(),
      value: 'visible',
    },
    {
      label: $LL.brainstormVisibilityLabelConcealed(),
      value: 'concealed',
    },
    {
      label: $LL.brainstormVisibilityLabelHidden(),
      value: 'hidden',
    },
  ];

  function createRetro(e) {
    e.preventDefault();
    let endpoint = `${apiPrefix}/users/${$user.id}/retros`;
    const body = {
      retroName,
      format: 'worked_improve_question',
      joinCode,
      facilitatorCode,
      maxVotes,
      brainstormVisibility,
    };

    if (selectedTeam !== '') {
      endpoint = `/api/teams/${selectedTeam}/users/${$user.id}/retros`;
    }

    xfetch(endpoint, { body })
      .then(res => res.json())
      .then(function ({ data }) {
        eventTag('create_retro', 'engagement', 'success', () => {
          router.route(`${appRoutes.retro}/${data.id}`);
        });
      })
      .catch(function (error) {
        if (Array.isArray(error)) {
          error[1].json().then(function (result) {
            notifications.danger(
              `${$LL.createRetroErrorMessage()} : ${result.error}`,
            );
          });
        } else {
          notifications.danger($LL.createRetroErrorMessage());
        }
        eventTag('create_retro', 'engagement', 'failure');
      });
  }

  function getTeams() {
    xfetch(`/api/users/${$user.id}/teams?limit=100`)
      .then(res => res.json())
      .then(function (result) {
        teams = result.data;
      })
      .catch(function () {
        notifications.danger($LL.getTeamsError());
      });
  }

  onMount(() => {
    if (!$user.id) {
      router.route(appRoutes.register);
    }
    getTeams();
  });
</script>

<form on:submit="{createRetro}" name="createRetro">
  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="retroName"
    >
      {$LL.retroName()}
    </label>
    <div class="control">
      <input
        name="retroName"
        bind:value="{retroName}"
        placeholder="{$LL.retroNamePlaceholder()}"
        class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2
            appearance-none
            rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
            focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500
            dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
        id="retroName"
        required
      />
    </div>
  </div>

  {#if apiPrefix === '/api'}
    <div class="mb-4">
      <label
        class="text-gray-700 dark:text-gray-400 text-sm font-bold inline-block mb-2"
        for="selectedTeam"
      >
        {$LL.associateTeam()}
        {#if !AppConfig.RequireTeams}{$LL.optional()}
        {/if}
      </label>
      <div class="relative">
        <select
          bind:value="{selectedTeam}"
          class="block appearance-none w-full border-2 border-gray-300 dark:border-gray-700
                text-gray-700 dark:text-gray-300 py-3 px-4 pe-8 rounded leading-tight
                focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 dark:bg-gray-900"
          id="selectedTeam"
          name="selectedTeam"
        >
          <option value="" disabled>{$LL.selectTeam()}</option>
          {#each teams as team}
            <option value="{team.id}">
              {team.name}
            </option>
          {/each}
        </select>
        <div
          class="pointer-events-none absolute inset-y-0 end-0 flex
                items-center px-2 text-gray-700 dark:text-gray-400"
        >
          <DownCarrotIcon />
        </div>
      </div>
    </div>
  {/if}

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="joinCode"
    >
      {$LL.joinCodeLabelOptional()}
    </label>
    <div class="control">
      <input
        name="joinCode"
        bind:value="{joinCode}"
        placeholder="{$LL.joinCodePlaceholder()}"
        class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
                focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
        id="joinCode"
      />
    </div>
  </div>

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="facilitatorCode"
    >
      {$LL.facilitatorCodeOptional()}
    </label>
    <div class="control">
      <input
        name="facilitatorCode"
        bind:value="{facilitatorCode}"
        placeholder="{$LL.facilitatorCodePlaceholder()}"
        class="bg-gray-100 dark:bg-gray-900 dark:focus:bg-gray-800 border-gray-200 dark:border-gray-600 border-2
            appearance-none
            rounded w-full py-2 px-3 text-gray-700 dark:text-gray-400 leading-tight
            focus:outline-none focus:bg-white focus:border-indigo-500 focus:caret-indigo-500
            dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
        id="facilitatorCode"
      />
    </div>
  </div>

  <div class="mb-4">
    <label
      class="block text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="maxVotes"
    >
      {$LL.retroMaxVotesPerUserLabel()}
    </label>
    <div class="control">
      <input
        name="maxVotes"
        bind:value="{maxVotes}"
        class="bg-gray-100 dark:bg-gray-900 border-gray-200 dark:border-gray-800 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 dark:text-gray-300 leading-tight
                focus:outline-none focus:bg-white dark:focus:bg-gray-700 focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
        id="maxVotes"
        type="number"
        min="1"
        max="10"
        required
      />
    </div>
  </div>

  <div class="mb-4">
    <label
      class="text-gray-700 dark:text-gray-400 text-sm font-bold mb-2"
      for="brainstormVisibility"
    >
      {$LL.brainstormPhaseFeedbackVisibility()}
    </label>
    <div class="relative">
      <select
        bind:value="{brainstormVisibility}"
        class="block appearance-none w-full border-2 border-gray-300 dark:border-gray-700
                text-gray-700 dark:text-gray-300 py-3 px-4 pe-8 rounded leading-tight
                focus:outline-none focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400 dark:bg-gray-900"
        id="brainstormVisibility"
        name="brainstormVisibility"
      >
        {#each brainstormVisibilityOptions as item}
          <option value="{item.value}">
            {item.label}
          </option>
        {/each}
      </select>
      <div
        class="pointer-events-none absolute inset-y-0 end-0 flex
                items-center px-2 text-gray-700 dark:text-gray-400"
      >
        <DownCarrotIcon />
      </div>
    </div>
  </div>

  <div class="text-right">
    <SolidButton type="submit">{$LL.createRetro()}</SolidButton>
  </div>
</form>
