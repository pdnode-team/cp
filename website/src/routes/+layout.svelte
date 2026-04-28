<script lang="ts">
	import './layout.css'

	let { children } = $props()
	let user: any = $state()

	import pb from '$lib/pocketbase'
	import { goto } from '$app/navigation'
	import { onMount } from 'svelte'

	function reloadLoginStatus() {
		if (pb.authStore.isValid) {
			user = pb.authStore.record
		} else {
			user = null
		}
		
	}
	pb.authStore.onChange((_, record) => {
		reloadLoginStatus()
	}, true)

	onMount(() => {
		reloadLoginStatus()
	})
</script>

<!-- 顶部导航栏 -->
<header
	class="navbar mx-auto max-w-7xl border-b border-base-200 bg-base-100 px-4 shadow-sm md:px-8"
>
	<div class="flex-1">
		<a href="/" class="text-xl font-bold tracking-tight transition-colors hover:text-primary">
			Pdnode CP Website <div class="badge badge-soft">v0.5.1</div>
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
				<ul class="dropdown-content menu z-1 mt-3 w-52 menu-sm rounded-box bg-base-300 p-2 shadow flex flex-col gap-2">
					<li>
						<button
							onclick={() => {
								pb.authStore.clear()
								window.location.reload()
							}}
							class="w-full text-left text-error">Logout</button
						>
					</li>
					<li>
						<button onclick={() => goto("/mystuff")}>My Stuff</button>
					</li>
				</ul>
			</div>
		{:else}
			<a href="/register" class="btn btn-outline btn-sm">Register</a>
			<a href="/login" class="btn btn-sm btn-primary">Login</a>
		{/if}
	</div>
</header>

<!-- 主体内容插槽 -->
<main class="mx-auto min-h-[calc(100vh-65px)] max-w-7xl border-base-200 bg-base-100 md:border-x">
	{@render children()}
</main>
