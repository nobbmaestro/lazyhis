# Source this in your ~/.zshrc
autoload -U add-zsh-hook

zmodload zsh/datetime 2>/dev/null

# If zsh-autosuggestions is installed, configure it to use LazyHis's search. If
# you'd like to override this, then add your config after the $(lazyhis init zsh)
# in your .zshrc
_zsh_autosuggest_strategy_lazyhis() {
	suggestion=$(lazyhis search --exit-code 0 --limit 1 -- "$@")
}

if [ -n "${ZSH_AUTOSUGGEST_STRATEGY:-}" ]; then
	ZSH_AUTOSUGGEST_STRATEGY=("lazyhis" "${ZSH_AUTOSUGGEST_STRATEGY[@]}")
else
	ZSH_AUTOSUGGEST_STRATEGY=("lazyhis")
fi

LAZYHIS_HISTORY_ID=""

_lazyhis_preexec() {
	local id
	id=$(lazyhis history add -- "$1")
	export LAZYHIS_HISTORY_ID="$id"
	__lazyhis_preexec_time=${EPOCHREALTIME-}
}

_lazyhis_precmd() {
	local EXIT="$?" __lazyhis_precmd_time=${EPOCHREALTIME-}

	[[ -z "${LAZYHIS_HISTORY_ID:-}" ]] && return

	local duration=""
	if [[ -n $__lazyhis_preexec_time && -n $__lazyhis_precmd_time ]]; then
		printf -v duration %.0f $(((__lazyhis_precmd_time - __lazyhis_preexec_time) * 1000))
	fi

	(
		lazyhis history edit \
			--exit-code $EXIT \
			${duration:+--duration=$duration} \
			-- $LAZYHIS_HISTORY_ID &
	) >/dev/null 2>&1

	export LAZYHIS_HISTORY_ID=""
}

_lazyhis_search() {
	emulate -L zsh
	zle -I

	local output
	output=$(lazyhis $* -- "$BUFFER" 3>&1 1>&2 2>&3)

	zle reset-prompt

	if [[ -n $output ]]; then
		RBUFFER=""
		LBUFFER=$output

		if [[ $LBUFFER == __lazyhis_accept__:* ]]; then
			LBUFFER=${LBUFFER#__lazyhis_accept__:}
			zle accept-line
		elif [[ $LBUFFER == __lazyhis_prefill__:* ]]; then
			LBUFFER=${LBUFFER#__lazyhis_prefill__:}
		fi
	fi
}
_lazyhis_search_vicmd() {
	_lazyhis_search
}
_lazyhis_search_viins() {
	_lazyhis_search
}

_lazyhis_up_search() {
	# Only trigger if the buffer is a single line
	if [[ ! $BUFFER == *$'\n'* ]]; then
		_lazyhis_search "$@"
	else
		zle up-line
	fi
}
_lazyhis_up_search_vicmd() {
	_lazyhis_up_search
}
_lazyhis_up_search_viins() {
	_lazyhis_up_search
}

add-zsh-hook preexec _lazyhis_preexec
add-zsh-hook precmd _lazyhis_precmd

zle -N lazyhis-search _lazyhis_search
zle -N lazyhis-search-vicmd _lazyhis_search_vicmd
zle -N lazyhis-search-viins _lazyhis_search_viins
zle -N lazyhis-up-search _lazyhis_up_search
zle -N lazyhis-up-search-vicmd _lazyhis_up_search_vicmd
zle -N lazyhis-up-search-viins _lazyhis_up_search_viins

bindkey -M emacs '^r' lazyhis-search
bindkey -M viins '^r' lazyhis-search-viins
bindkey -M vicmd '^r' lazyhis-search-vicmd

bindkey -M emacs '^[[A' lazyhis-up-search
bindkey -M vicmd '^[[A' lazyhis-up-search-vicmd
bindkey -M viins '^[[A' lazyhis-up-search-viins
bindkey -M emacs '^[OA' lazyhis-up-search
bindkey -M vicmd '^[OA' lazyhis-up-search-vicmd
bindkey -M viins '^[OA' lazyhis-up-search-viins
bindkey -M vicmd 'k' lazyhis-up-search-vicmd
