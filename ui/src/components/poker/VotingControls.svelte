<script lang="ts">
  import SolidButton from '../SolidButton.svelte';
  import LL from '../../i18n/i18n-svelte';
  import ChevronDown from '../icons/ChevronDown.svelte';
  import { AppConfig } from '../../config';

  export let sendSocketEvent = () => {};
  export let eventTag;
  export let planId = '';
  export let points = [];
  export let votingLocked = true;
  export let highestVote = '';

  $: planPoints = highestVote;

  const endPlanVoting = () => {
    sendSocketEvent('end_voting', planId);
    eventTag('vote_end', 'battle', '');
  };

  const skipPlan = () => {
    sendSocketEvent('skip_plan', planId);
    eventTag('plan_skip', 'battle', '');
  };

  const restartVoting = () => {
    sendSocketEvent('activate_plan', planId);
    eventTag('plan_restart_vote', 'battle', '');
  };

  function handleSubmit(event) {
    event.preventDefault();

    sendSocketEvent(
      'finalize_plan',
      JSON.stringify({
        planId,
        planPoints,
      }),
    );
    eventTag('plan_finalize', 'battle', planPoints);

    planPoints = '';
  }
</script>

{#if planId != ''}
  <div class="p-4" data-testId="votingControls">
    <SolidButton
      color="blue"
      additionalClasses="mb-2 w-full"
      onClick="{skipPlan}"
      testid="voting-skip"
    >
      {$LL.planSkip({ friendly: AppConfig.FriendlyUIVerbs })}
    </SolidButton>
    {#if !votingLocked}
      <SolidButton
        additionalClasses="w-full"
        onClick="{endPlanVoting}"
        testid="voting-finish"
      >
        {$LL.votingFinish()}
      </SolidButton>
    {:else}
      <SolidButton
        color="blue"
        additionalClasses="mb-2 w-full"
        onClick="{restartVoting}"
        testid="voting-restart"
      >
        {$LL.votingRestart()}
      </SolidButton>
      <form on:submit="{handleSubmit}" name="savePlanPoints">
        <legend
          class="text-xl mb-2 font-semibold leading-tight dark:text-gray-300"
        >
          {$LL.finalPoints()}
        </legend>
        <div class="flex -mx-2">
          <div class="w-1/2 px-2">
            <div class="relative">
              <select
                name="planPoints"
                bind:value="{planPoints}"
                required
                class="block appearance-none w-full border-2 dark:bg-gray-900
                                border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-300 py-3 px-4 pe-8
                                rounded leading-tight focus:outline-none
                                focus:border-indigo-500 focus:caret-indigo-500 dark:focus:border-yellow-400 dark:focus:caret-yellow-400"
              >
                <option value="" disabled>
                  {$LL.points()}
                </option>
                >
                {#each points as point}
                  <option value="{point}">{point}</option>
                {/each}
              </select>
              <div
                class="pointer-events-none absolute inset-y-0
                                end-0 flex items-center px-2 text-gray-700 dark:text-gray-300"
              >
                <ChevronDown />
              </div>
            </div>
          </div>
          <div class="w-1/2 text-right px-2">
            <SolidButton
              additionalClasses="w-full h-full"
              type="submit"
              testid="voting-save"
            >
              {$LL.save()}
            </SolidButton>
          </div>
        </div>
      </form>
    {/if}
  </div>
{/if}
