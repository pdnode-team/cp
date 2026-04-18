<script lang="ts">
	import './layout.css';

	export const ssr = false;
	let { children } = $props();
	let user: any = $state();

	import pb from '$lib/pocketbase';
	import { onMount } from 'svelte';

	onMount(() => {
		user = pb.authStore.record;
	});
</script>

<!-- 顶部导航栏 -->
<header
	class="navbar mx-auto max-w-7xl border-b border-base-200 bg-base-100 px-4 shadow-sm md:px-8"
>
	<div class="flex-1">
		<a href="/" class="text-xl font-bold tracking-tight transition-colors hover:text-primary">
			Pdnode CP Website
		</a>
	</div>
	<div class="flex-none gap-4">
		<a href="/about" class="btn text-base-content/80 btn-ghost btn-sm">About</a>

		{#if user}
			<a href="/explore" class="btn text-base-content/80 btn-ghost btn-sm">Explore</a>
			<div class="dropdown dropdown-end">
				<button class="placeholder btn avatar btn-circle btn-ghost">
					<div
						class="flex w-10 items-center justify-center rounded-full bg-primary text-primary-content"
					>
						<span class="text-lg leading-none font-bold">
							{user.name?.trim().at(0) ?? '?'}
						</span>
					</div>
				</button>
				<ul class="dropdown-content menu z-1 mt-3 w-52 menu-sm rounded-box bg-base-300 p-2 shadow">
					<li>
						<button
							onclick={() => {
								pb.authStore.clear();
								window.location.reload();
							}}
							class="w-full text-left text-error">Logout</button
						>
					</li>
				</ul>
			</div>
		{:else}
			<a href="/register" class="btn btn-outline btn-sm">Register</a>
			<a href="/login" class="btn btn-sm btn-primary">Login</a>
		{/if}
	</div>
</header>

<!-- 全局提示消息 (Toast)
{#if flashMessages.error || flashMessages.success}
  <div class="toast toast-top toast-center z-50 mt-16">
    {#if flashMessages.error}
      <div class="alert alert-error shadow-lg">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
        <span>{flashMessages.error}</span>
      </div>
    {/if}
    {#if flashMessages.success}
      <div class="alert alert-success shadow-lg">
        <svg xmlns="http://www.w3.org/2000/svg" class="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
        <span>{flashMessages.success}</span>
      </div>
    {/if}
  </div>
{/if} -->

<!-- 主体内容插槽 -->
<main class="mx-auto min-h-[calc(100vh-65px)] max-w-7xl border-base-200 bg-base-100 md:border-x">
	{@render children()}
</main>
